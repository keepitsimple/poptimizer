Оптимизация долгосрочного портфеля акций
========================================
[![image](https://github.com/WLM1ke/poptimizer/workflows/tests/badge.svg)](https://github.com/WLM1ke/poptimizer/actions)
[![image](https://codecov.io/gh/WLM1ke/poptimizer/branch/master/graph/badge.svg)](https://codecov.io/gh/WLM1ke/poptimizer)

О проекте
---------

Занимаюсь инвестициями с 2008 года. Целью проекта является изучение
программирования и автоматизация процесса управления портфелем акций.

Используемый подход не предполагает баснословных доходностей, а нацелен
на получение результата чуть лучше рынка при рисках чуть меньше рынка
при относительно небольшом обороте. Портфель ценных бумаг должен быть
достаточно сбалансированным, чтобы его нестрашно было оставить без
наблюдения на продолжительное время.

Большинство частных инвесторов стремиться к быстрому обогащению и,
согласно известному афоризму Баффета, \"мало кто хочет разбогатеть
медленно\", поэтому проект является открытым. Стараюсь по возможности
исправлять ошибки, выявленные другими пользователями, и буду рад любой
помощи от более опытных программистов. Особенно приветствуются вопросы и
предложения по усовершенствованию содержательной части подхода к
управлению портфелем.

Проект находится в стадии развития и постоянно модифицируется (не всегда
удачно), поэтому может быть использован на свой страх и риск.

Основные особенности
--------------------

### Оптимизация портфеля

- Базируется на [Modern portfolio
    theory](https://en.wikipedia.org/wiki/Modern_portfolio_theory)
- При построении портфеля учитывается более 230 акций (включая
    иностранные) и ETF, обращающихся на MOEX
- Используется ансамбль моделей для оценки неточности предсказаний
    ожидаемых доходностей и рисков отдельных активов
- Используется робастная инкрементальная оптимизация на основе расчета
    улучшения метрик портфеля в результате торговли с учетом неточности 
    имеющихся прогнозов вместо классической
    mean-variance оптимизации
- Применяется [поправка
    Бонферрони](https://en.wikipedia.org/wiki/Bonferroni_correction) на
    множественное тестирование с учетом большое количества анализируемых
    активов

### Прогнозирование параметров активов

-   Используются нейронные сети на основе архитектуры
    [WaveNet](https://arxiv.org/abs/1609.03499) с большим receptive
    field для анализа длинных последовательностей котировок
-   Осуществляется совместное прогнозирование ожидаемой доходности и ее
    дисперсии с помощью подходов, базирующихся на [GluonTS:
    Probabilistic Time Series Models in
    Python](https://arxiv.org/abs/1906.05264)
-   Для моделирования толстых хвостов в распределениях доходностей
    применяются смеси логнормальных распределений
-   Используются устойчивые оценки исторических корреляционных матриц
    для большого числа активов с помощью сжатия
    [Ledoit-Wolf](http://www.ledoit.net/honey.pdf)

### Формирование ансамбля моделей

-   Осуществляется выбор моделей из многомерного пространства
    гиперпараметров сетей, их оптимизаторов и комбинаций признаков
-   Для исследования пространства применяются подходы алгоритма
    [Имитации
    отжига](https://en.wikipedia.org/wiki/Simulated_annealing)
-   Для масштабирования локальной области поиска и кодирования
    гиперпараметров используются принципы [дифференциальной
    эволюции](https://en.wikipedia.org/wiki/Differential_evolution)
-   Для выбора моделей в локальной области применяется распределение
    [Коши](https://en.wikipedia.org/wiki/Cauchy_distribution) для
    осуществления редких не локальных прыжков в пространстве
    гиперпараметров
-   При отборе претендентов в ансамбль осуществляется [последовательное
    тестирование](https://en.wikipedia.org/wiki/Sequential_analysis#Alpha_spending_functions)
    с соответствующими корректировками [уровней
    значимости](https://arxiv.org/abs/1906.09712)

### Источники данных

-   Реализована загрузка котировок всех акций (включая иностранные) и
    ETF, обращающихся на MOEX
-   Поддерживается в актуальном состоянии база данных дивидендов с 2015г
    по включенным в анализ акциям
-   Реализована возможность сверки базы данных дивидендов с информацией
    на сайтах:

> -   [www.dohod.ru](https://www.dohod.ru/ik/analytics/dividend)
> -   [www.conomy.ru](https://www.conomy.ru/dates-close/dates-close2)
> -   [bcs-express.ru](https://bcs-express.ru/dividednyj-kalendar)
> -   [www.smart-lab.ru](https://smart-lab.ru/dividends/index/order_by_yield/desc/)
> -   [закрытияреестров.рф](https://закрытияреестров.рф/)
> -   [finrange.com](https://finrange.com/)
> -   [investmint.ru](https://investmint.ru/)
> -   [www.nasdaq.com](https://www.nasdaq.com/)
> -   [www.streetinsider.com](https://www.streetinsider.com/)

Направления дальнейшего развития
--------------------------------

- Реализация сервиса на Go для загрузки всей необходимой информации
- Применение нелинейного сжатия Ledoit-Wolf для оценки корреляции
    активов
- Рефакторинг кода на основе
    [DDD](https://en.wikipedia.org/wiki/Domain-driven_design),
    [MyPy](http://mypy.readthedocs.org/en/latest/) и
    [wemake](https://wemake-python-stylegui.de/en/latest/)
- Использование архитектур на основе
    [трансформеров](https://en.wikipedia.org/wiki/Transformer_(machine_learning_model))
    вместо WaveNet

FAQ
---

### Какие инструменты нужны для запуска программы?

Последняя версия MongoDB, MongoDB Database Tools, Python и все зависимости из [requirements.txt](https://github.com/WLM1ke/poptimizer/blob/18756e8bdbfcac3ebd7ba241f86b25bcb27cc22f/requirements.txt).

### Как запускать программу?

Запуск реализован через CLI:

`python3 -m poptimizer`

После этого можно посмотреть перечень команд и help к ним, а дальше самому разбираться в коде.
Основные команды для запуска описаны в файле [\_\_main\_\_.py](https://github.com/WLM1ke/poptimizer/blob/18756e8bdbfcac3ebd7ba241f86b25bcb27cc22f/poptimizer/__main__.py).
Сначала необходимо запустить функцию `evolve` для обучения моделей. После этого можно запустить `optimize` для 
оптимизации портфеля.

### Что за код в папке opt

Пока находящаяся в разработке реализация на Go с web-интерфейсом. Так как формат базы данных не совместим, то 
следует запускать на отдельной MongoDB.

### У меня появилось сообщение ДАННЫЕ ПО ДИВИДЕНДАМ ТРЕБУЮТ ОБНОВЛЕНИЯ - что делать?

Вся необходимая для работы программы информация обновляется автоматически, кроме данных по дивидендам, которые
необходимо обновлять в ручную в базе данных `source`. После ввода информации необходимо запустить команду `dividends`
для дополнительной сверки с разными внешними источниками данных и перемещения информации в основную рабочую базу `data`.

### Есть ли у программы какие-нибудь настройки?

Настройки описаны в файле [config.template](https://github.com/WLM1ke/poptimizer/blob/18756e8bdbfcac3ebd7ba241f86b25bcb27cc22f/config/config.template).
При отсутствии файла конфигурации будут использоваться значения по умолчанию. 

### Как ввести свой портфель?

Пример заполнения файла с портфеля с базовым набором бумаг содержится в файле [base.yaml](https://github.com/WLM1ke/poptimizer/blob/18756e8bdbfcac3ebd7ba241f86b25bcb27cc22f/portfolio/base.yaml).
В этом каталоге можно хранить множество файлов (например по отдельным брокерским счетам), информация из них будет 
объединяться в единый портфель. Для работы оптимизатора необходимо наличие хотя бы одной не нулевой позиции по бумагам.

### У меня в портфеле не так много бумаг - нулевые значения по множеству позиций как-нибудь влияют на эволюцию?

Для эволюции принципиален перечень бумаг, а не их количество. 

### У меня в портфеле не так много бумаг - можно оставить только свои?

При желании можно сократить, количество бумаг, но их должно быть не меньше половины из базового набора, чтобы получалось 
достаточно большое количество обучающих примеров для тренировки моделей.

### В моем портфеле есть бумаги, отсутствующие в базовом наборе - можно ли их добавить?

Можно добавить любые акции, включая иностранные, и ETF, обращающиеся на MOEX. Для корректной работы так же может 
потребоваться дополнить базу по дивидендам, если они выплачивались с 2015 года.

### Я добавил новую бумагу, запустил эволюцию и получил ошибки "Удаляю - Слишком большая длинна истории" для всех моделей. В результате вся популяция погибла. Как избежать такой проблемы?

Для обучения модели нужно, чтобы по каждому тикеру был хотя бы один обучающий пример и некоторое количество тестирующий 
примеров (зависит от количества поколений). Для этого нужна минимальная история котировок порядка history_days + 
forecast_days. Если истории не хватает, то модель удаляется и выводится указанное сообщение. Когда добавляется новый 
тикер с очень короткой историей, теоретически может погибнуть вся популяция.

Перед добавлением тикеров рекомендуется вызывать метод add_tickers для старого портфеля. В результате будут выведены 
тикеры (торгуемые на бирже, но еще не включенные в портфель), у которых с некоторым запасом есть минимально необходимая 
история. Дополнительно отбрасываются совсем малоликвидные с точки зрения размера вашего портфеля бумаги. Если он 
большой, то отбрасывается больше тикеров. Выбранные бумаги упорядочиваются по возрастанию корреляции - тикеры в начале 
списка более удачны с точки зрения потенциала снижения рисков.

Бумаги из этого списка относительно безопасно включать в анализ - некоторые модели с очень большим параметром 
history_days могут погибнуть, но основной костяк будет работать.

### Можно ли остановить эволюцию (прервать через Ctrl-C), чтобы, например, перезапустить компьютер?

Можно останавливать в любое время, но желательно в момент обучения, а не в момент тестирования - если идет тестирование 
лучше дождаться его завершения и начала обучения следующей модели, но это не сильно критично.

### Верно ли, что эволюция должна быть постоянно запущена?

При настройках по умолчанию предполагается, что эволюция будет работать постоянно, но в принципе достаточно, чтобы она 
запускалась регулярно и успевала проходить полный круг для существующих моделей после появления новых данных.

При желании в конфиге можно указать часы между которыми будет осуществляться тренировка (START_EVOLVE_HOUR и 
STOP_EVOLVE_HOUR), например, ночью. 

### Я поменял перечень тикеров и при запуске оптимизации стали выпадать ошибки. Как правильно поменять перечень тикеров?

Для эволюции не важно, когда вы изменили перечень тикеров. Она всегда берет текущий список и тренирует новые модели 
под него, а модели уже протестированные для текущего дня будут переучены после поступления новых данных по котировкам на
следующий день.

Для оптимизатора подходят только модели с соответствующим набором тикеров, поэтому если вы чего-то поменяли сразу после 
этого не будет моделей, подходящих для оптимизации, что будет вызывать ошибки. Как только эволюция обучит хотя бы две 
новые модели, оптимизатор сможет работать без ошибок, но этого будет недостаточно для адекватной оптимизации. 

Лучше всего менять тикеры после того, как вы поторговали и больше не собираетесь торговать в этот день - тогда у вас 
будет время на обучение новых моделей под новый набор тикеров, а после 0:45 появятся данные и начнется переобучение 
старых моделей под новые тикеры. Утром к началу торгов у вас будет какое-то разумное количество моделей 
актуализированных под новый набор бумаг. 

Еще лучше менять перечень тикеров вечером в пятницу, особенно если у вас популяция в несколько сотен моделей, так как на 
ее переобучение вполне могут уйти все выходные.

### Нужно ли очищать базу с моделями после изменения тикеров?

Не нужно, более того, лучше ее сохранять. Модели, которые работали хорошо раньше, обычно хорошо работают на других 
тикерах, просто им нужно переобучиться. Эволюция начинается с заведом плохих моделей, поэтому регулярно сбрасывая базу 
вы будете постоянно использовать слабые модели и долго ждать их эволюционного улучшения.

### Нужно ли запускать оптимизацию портфеля каждый день?

По умолчанию предполагается, что торговля будет осуществляться каждый торговый день. Если вы торгуете реже, поменяйте 
значение параметра TRADING_INTERVAL, например на 5 для торговли раз в неделю. В большинстве случаев после завершения 
первичного процесса диверсификации, торговые сигналы будут появляться далеко не каждый день и будут затрагивать 
небольшую долю портфеля, поэтому даже при дефолтных настройках, торговля не будет отнимать много усилий.  

### Как рекомендуется исполнять торговые рекомендации?

Оптимизатор использует линейную аппроксимацию для построения рекомендаций. Предполагается, что вы должны совершать 
небольшие сделки (~1% от стоимости портфеля), потом повторно запускать оптимизатор для уточнения рекомендаций и т.д. 
пока не останется рекомендаций на продажу. В результате останется одна рекомендация (PENDING) - это лучшая рекомендация 
на покупку, чтобы было ясно, куда направлять кэш, если он остался после последних операций или внесен на счет. 

Рекомендации обычно сохраняются в течении продолжительного времени, поэтому нет большой необходимости пытаться исполнить 
их за один день, особенно если у вас большой по объему портфель или операции затрагивают малоликвидные бумаги - торгуйте
в удобном для вас темпе.

Если у вас на счете по каким-то причинам много кэша (много внесли или случайно продали на большую сумму) сначала
покупайте маленькими порциями, пока не снизите его количество.
 
### Что отражают LOWER и UPPER в разделе оптимизация портфеля?

Нижняя и верхняя граница доверительного интервала влияния указанной бумаги на качество портфеля. Если верхняя граница 
одной бумаги ниже нижней границы второй бумаги, то целесообразно сокращать позицию по первой бумаге и наращивать позицию 
по второй. При выдаче рекомендаций дополнительно учитывается, что зазор между границами должен покрывать транзакционные 
издержки.

Особые благодарности
--------------------

-   [Evgeny Pogrebnyak](https://github.com/epogrebnyak) за помощь в
    освоении Python
-   [RomaKoks](https://github.com/RomaKoks) за полезные советы по
    автоматизации некоторых этапов работы программы и исправлению ошибок
-   [AlexQww](https://github.com/AlexQww) за содержательные обсуждения
    подходов к управлению портфелем, которые стали катализатором
    множества изменений в программе
