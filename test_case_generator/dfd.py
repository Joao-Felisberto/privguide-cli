from typing import List, Self

from util import Cloneable


class ProcessedCategorization(Cloneable):
    def __init__(
            self,
            level_0: str,
            level_1: str,
    ):
        self.level_0: str = level_0
        self.level_1: str = level_1

    @staticmethod
    def default() -> 'ProcessedCategorization':
        return ProcessedCategorization("logical", "boolean")


class DataFormat(Cloneable):
    def __init__(
            self,
            format_type: str,
            level_0: str,
            level_1: str,
    ):
        self.format_type: str = format_type
        self.level_0: str = level_0
        self.level_1: str = level_1

    @staticmethod
    def default() -> 'DataFormat':
        return DataFormat("psdc:plain", "plain text", "business")


class Data(Cloneable):
    def __init__(
            self,
            id_: str,
            format: DataFormat,
            data_processed: ProcessedCategorization,
    ):
        self.id_: str = id_
        self.format: DataFormat = format
        self.data_processed: ProcessedCategorization = data_processed

#     @staticmethod
#     def default() -> Self:
#         return Self("", DataFormat.default())

class ExternalEntity(Cloneable):
    def __init__(
            self,
            id_: str,
            consumes: List[Data],
            produces: List[Data],
            location: List[str],
            environment: List[str],
            categories: List[str],
            age: str | None,
            produces_public_information: bool,
            safeguards: List[str] | str,
            options: List[str] | str,
    ):
        self.id_: str = id_
        self.consumes: List[Data] = consumes
        self.produces: List[Data] = produces
        self.location: List[str] = location
        self.environment: List[str] = environment
        self.categories: List[str] = categories
        self.age: str | None = age
        self.produces_public_information: bool = produces_public_information
        self.safeguards: List[str] | str = safeguards
        self.options: List[str] | str = options


class DataStored(Cloneable):
    def __init__(
            self,
            type_: str,
            storage_period: str,
            create: str,
            read: str,
            update: str,
            delete: str,
            format_: DataFormat,
    ):
        self.type_: str = type_
        self.storage_period: str = storage_period
        self.create: str = create
        self.read: str = read
        self.update: str = update
        self.delete: str = delete
        self.format_: DataFormat = format_


class DataFlow(Cloneable):
    def __init__(
            self,
            id_: str,
            from_: str,
            to: str,
            data: List[Data],
            encryption: str,
            periodicity: str,
            amount_of_data_per_period: int,
            certification: List[str],
            safeguards: List[str],
    ):
        self.id_: str = id_
        self.from_: str = from_
        self.to: str = to
        self.data: List[Data] = data
        self.encryption: str = encryption
        self.periodicity: str = periodicity
        self.amount_of_data_per_period: int = amount_of_data_per_period
        self.certification: List[str] = certification
        self.safeguards: List[str] = safeguards


class DataStore(Cloneable):
    def __init__(
            self,
            id_: str,
            data_stored: List[DataStored],
            location: List[str],
            environment: List[str],
            certification: List[str],
            safeguards: List[str],
    ):
        self.id_: str = id_
        self.data_stored: List[DataStored] = data_stored
        self.location: List[str] = location
        self.environment: List[str] = environment
        self.certification: List[str] = certification
        self.safeguards: List[str] = safeguards


class Process(Cloneable):
    def __init__(
            self,
            id_: str,
            consumes: List[Data],
            produces: List[Data],
            location: List[str],
            environment: List[str],
            purposes: List[str],
            certification: List[str],
            safeguards: List[str],
    ):
        self.id_: str = id_
        self.consumes: List[Data] = consumes
        self.produces: List[Data] = produces
        self.location: List[str] = location
        self.environment: List[str] = environment
        self.purposes: List[str] = purposes
        self.certification: List[str] = certification
        self.safeguards: List[str] = safeguards


class DataType(Cloneable):
    def __init__(
            self,
            id_: str,
            aggregates: List[str],
            validity: str,
            categories: List[str],
    ):
        self.id_: str = id_
        self.aggregates: List[str] = aggregates
        self.validity: str = validity
        self.categories: List[str] = categories


class DFD(Cloneable):
    def __init__(
            self,
            data_types: List[DataType],
            external_entities: List[ExternalEntity],
            processes: List[Process],
            data_stores: List[DataStore],
            data_flows: List[DataFlow],
    ):
        self.data_types: List[DataType] = data_types
        self.external_entities: List[ExternalEntity] = external_entities
        self.processes: List[Process] = processes
        self.data_stores: List[DataStore] = data_stores
        self.data_flows: List[DataFlow] = data_flows
