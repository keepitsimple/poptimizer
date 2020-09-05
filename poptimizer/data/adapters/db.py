"""Реализации сессий доступа к базе данных."""
import asyncio
import logging
from typing import Iterable, Optional, Tuple

import pandas as pd

from poptimizer.data.config import resources
from poptimizer.data.ports import base, outer

# База данных и коллекция для одиночный
DB = "data_new"
MISC = "misc"


def _collection_and_name(table: outer.TableTuple) -> Tuple[str, str]:
    """Формирует название коллекции и имя документа."""
    collection: str = table.group
    name = table.name
    if collection == name:
        collection = MISC
    return collection, name


class MongoDBSession(outer.AbstractDBSession):
    """Реализация сессии с хранением в MongoDB.

    При совпадении id и группы данные записываются в специальную коллекцию, в ином случае в коллекцию
    группы.
    """

    def __init__(self) -> None:
        """Получает ссылку на базу данных."""
        self._logger = logging.getLogger(self.__class__.__name__)
        client = resources.get_mongo_client()
        self._db = client[DB]

    async def get(self, table_name: base.TableName) -> Optional[outer.TableTuple]:
        """Извлекает документ из коллекции."""
        group, name = table_name
        collection: str = group
        if collection == name:
            collection = MISC
        doc = await self._db[collection].find_one({"_id": name})

        if doc is None:
            return None

        df = pd.DataFrame(**doc["data"])
        return outer.TableTuple(group=group, name=name, df=df, timestamp=doc["timestamp"])

    async def commit(self, tables_vars: Iterable[outer.TableTuple]) -> None:
        """Записывает данные в MongoDB."""
        aws = []

        for table in tables_vars:
            collection, name = _collection_and_name(table)
            self._logger.info(f"Сохранение {collection}.{name}")

            aw_update = self._db[collection].replace_one(
                filter={"_id": name},
                replacement=dict(_id=name, data=table.df.to_dict("split"), timestamp=table.timestamp),
                upsert=True,
            )
            aws.append(aw_update)

        await asyncio.gather(*aws)
