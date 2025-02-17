"""Эволюция параметров модели."""
import datetime
import logging
import operator
from typing import Optional

import numpy as np

from poptimizer import config
from poptimizer.data.views import listing
from poptimizer.dl import ModelError
from poptimizer.evolve import population, seq
from poptimizer.portfolio.portfolio import load_tickers


class Evolution:  # noqa: WPS214
    """Эволюция параметров модели.

    Эволюция состоит из бесконечного создания организмов и сравнения их характеристик с медианными значениями по
    популяции. Сравнение осуществляется с помощью последовательного теста для медиан, который учитывает изменение
    значимости тестов при множественном тестировании по мере появления данных за очередной период времени. Дополнительно
    осуществляется коррекция на множественное тестирование на разницу llh и доходности.
    """

    def __init__(self):
        """Инициализирует необходимые параметры."""
        self._tickers = None
        self._end = None
        self._logger = logging.getLogger()
        self._scale = max(1, population.count())

    def evolve(self) -> None:
        """Осуществляет эволюции.

        При необходимости создается начальная популяция из случайных организмов по умолчанию.
        """
        step = 0
        org = None

        while _check_time_range():
            step = self._step_setup(step)

            date = self._end.date()
            self._logger.info(f"***{date}: Шаг эволюции — {step}***")
            population.print_stat()
            self._logger.info(f"Scale - {self._scale}\n")

            if org is None:
                org = population.get_next_one(self._end) or population.get_next_one(None)

            org = self._step(org)

    def _step_setup(
        self,
        step: int,
    ) -> int:
        self._setup()

        d_min, d_max = population.min_max_date()
        if self._tickers is None:
            self._tickers = load_tickers()
            self._end = d_max or listing.all_history_date(self._tickers)[-1]

        dates = listing.all_history_date(self._tickers, start=self._end)
        if (d_min != self._end) or (len(dates) == 1):
            return step + 1

        self._end = dates[1]

        return 1

    def _setup(self) -> None:
        if population.count() == 0:
            while population.count() < seq.minimum_bounding_n(config.P_VALUE / max(1, population.count())):
                self._logger.info("Создается базовый организм:")
                org = population.create_new_organism()
                self._logger.info(f"{org}\n")

            self._scale = max(1, population.count())

    def _maybe_clear(self, org: population.Organism) -> population.Organism:

        if (org.date == self._end) and (0 < org.scores < max(self._scale, _n_test())):
            org.clear()

        if (org.date != self._end) and (0 < org.scores < max(self._scale, _n_test()) - 1):
            org.clear()

        return org

    def _step(self, hunter: population.Organism) -> Optional[population.Organism]:
        """Один шаг эволюции."""
        skip = True

        if not hunter.scores or hunter.date == self._end:
            skip = False
            self._maybe_clear(hunter)

        label = ""
        if not hunter.scores:
            label = " - новый организм"

        self._logger.info(f"Родитель{label}:")
        if (margin := self._eval_organism(hunter)) is None:
            return None
        if margin[0] < 0:
            return None
        if skip:
            return None
        if margin[0] - margin[1] < 0:
            self._logger.info("Медленный не размножается...\n")

            return None

        prey = hunter.make_child(1 / max(1, self._scale))

        self._logger.info(f"Потомок:")
        if (margin := self._eval_organism(prey)) is None:
            return None
        if margin[0] < 0:
            self._scale += 1

            return None

        self._scale = max(1, self._scale - 1)

        if margin[0] - margin[1] < 0:
            return None

        return None

    def _eval_organism(self, organism: population.Organism) -> Optional[tuple[float, float]]:
        """Оценка организмов.

        - Если организм уже оценен для данной даты, то он не оценивается.
        - Если организм старый, то оценивается один раз.
        - Если организм новый, то он оценивается для определенного количества дат из истории.
        """
        try:
            self._logger.info(f"{organism}\n")
        except AttributeError as err:
            organism.die()
            self._logger.error(f"Удаляю - {err}\n")

            return None

        if organism.date == self._end:
            return self._get_margin(organism)

        dates = [self._end]
        if not organism.llh:
            dates = listing.all_history_date(self._tickers, end=self._end)
            dates = dates[-_n_test(organism.scores) :].tolist()

        for date in dates:
            try:
                organism.evaluate_fitness(self._tickers, date)
            except (ModelError, AttributeError) as error:
                organism.die()
                self._logger.error(f"Удаляю - {error}\n")

                return None

        return self._get_margin(organism)

    def _get_margin(self, org: population.Organism) -> tuple[float, float]:
        """Используется тестирование разницы llh и ret против самого старого организма.

        Используются тесты для связанных выборок, поэтому предварительно происходит выравнивание по
        датам и отбрасывание значений не имеющих пары (возможно первое значение и хвост из старых
        значений более старого организма).
        """
        margin = np.inf

        names = {"llh": "LLH", "ir": "RET"}

        for metric in ("llh", "ir"):
            median, upper, maximum = _select_worst_bound(
                candidate={"date": org.date, "llh": org.llh, "ir": org.ir},
                metric=metric,
            )

            self._logger.info(
                " ".join(
                    [
                        f"{names[metric]} worst difference:",
                        f"median - {median:0.4f},",
                        f"upper - {upper:0.4f},",
                        f"max - {maximum:0.4f}",
                    ],
                ),
            )

            valid = upper != median
            margin = min(margin, valid and (upper / (upper - median)))

        if margin == np.inf:
            margin = 0

        time_delta = _time_delta(org)

        self._logger.info(f"Margin - {margin:.2%}, Time excess - {time_delta:.2%}\n")  # noqa: WPS221

        if margin < 0:
            org.die()
            self._logger.info("Исключен из популяции...\n")

        return margin, time_delta


def _n_test(scores: int = -1) -> int:
    return max(population.count(), scores + 1)


def _time_delta(org):
    """Штраф за время, если организм медленнее медианного в популяции."""
    median = np.median([doc["timer"] for doc in population.get_metrics()])

    return max((org.timer / median - 1), 0)


def _check_time_range() -> bool:
    hour = datetime.datetime.today().hour

    if config.START_EVOLVE_HOUR == config.STOP_EVOLVE_HOUR:
        return True

    if config.START_EVOLVE_HOUR < config.STOP_EVOLVE_HOUR:
        return config.START_EVOLVE_HOUR <= hour < config.STOP_EVOLVE_HOUR

    before_midnight = config.START_EVOLVE_HOUR <= hour
    after_midnight = hour < config.STOP_EVOLVE_HOUR

    return before_midnight or after_midnight


def _select_worst_bound(candidate: dict, metric: str) -> tuple[float, float, float]:
    """Выбирает минимальное значение верхней границы доверительного интервала.

    Если данный организм не уступает целевому организму, то верхняя граница будет положительной.
    """

    diff = _aligned_diff(candidate, metric)

    bounds = map(
        lambda size: _test_diff(diff[:size]),
        range(1, len(diff) + 1),
    )

    return min(
        bounds,
        key=lambda bound: bound[1] or np.inf,
    )


def _aligned_diff(candidate: dict, metric: str) -> list[float]:
    comp = []

    for base in population.get_metrics():
        metrics = base[metric]

        if base["date"] < candidate["date"]:
            metrics = [np.nan] + metrics

        scores = len(candidate[metric])

        metrics = metrics[:scores]
        metrics = metrics + [np.nan] * (scores - len(metrics))

        comp.append(metrics)

    comp = np.nanmedian(np.array(comp), axis=0)

    return list(map(operator.sub, candidate[metric], comp))[::-1]


def _test_diff(diff: list[float]) -> tuple[float, float, float]:
    """Последовательный тест на медианную разницу с учетом множественного тестирования.

    Тестирование одностороннее, поэтому p-value нужно умножить на 2, но проводится 2 раза.
    """
    _, upper = seq.median_conf_bound(diff, config.P_VALUE / population.count())

    return float(np.median(diff)), upper, np.max(diff)
