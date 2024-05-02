from typing import List

from util import Cloneable


class ExternalEntity(Cloneable):
    def __init__(
            self,
            id_: str,
            consumes: [str],
            produces: [str],
            location: [str],
            environment: [str],
            categories: [str],
            age: str | None,
            produces_public_information: bool,
            safeguards: List[str] | str,
            options: List[str] | str,
    ):
        self.id_: str = id_
        self.consumes: [str] = consumes
        self.produces: [str] = produces
        self.location: [str] = location
        self.environment: [str] = environment
        self.categories: [str] = categories
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
    ):
        self.type_: str = type_
        self.storage_period: str = storage_period
        self.create: str = create
        self.read: str = read
        self.update: str = update
        self.delete: str = delete


class DataFlow(Cloneable):
    def __init__(
            self,
            id_: str,
            from_: str,
            to: str,
            data: [str],
            encryption: str,
            periodicity: str,
            amount_of_data_per_period: int,
            certification: [str],
            safeguards: [str],
    ):
        self.id_: str = id_
        self.from_: str = from_
        self.to: str = to
        self.data: [str] = data
        self.encryption: str = encryption
        self.periodicity: str = periodicity
        self.amount_of_data_per_period: int = amount_of_data_per_period
        self.certification: [str] = certification
        self.safeguards: [str] = safeguards


class DataStore(Cloneable):
    def __init__(
            self,
            id_: str,
            data_stored: [DataStored],
            location: [str],
            environment: [str],
            certification: [str],
            safeguards: [str],
    ):
        self.id_: str = id_
        self.data_stored: [DataStored] = data_stored
        self.location: [str] = location
        self.environment: [str] = environment
        self.certification: [str] = certification
        self.safeguards: [str] = safeguards


class Process(Cloneable):
    def __init__(
            self,
            id_: str,
            consumes: [str],
            produces: [str],
            location: [str],
            environment: [str],
            purposes: [str],
            certification: [str],
            safeguards: [str],
    ):
        self.id_: str = id_
        self.consumes: [str] = consumes
        self.produces: [str] = produces
        self.location: [str] = location
        self.environment: [str] = environment
        self.purposes: [str] = purposes
        self.certification: [str] = certification
        self.safeguards: [str] = safeguards


class DataType(Cloneable):
    def __init__(
            self,
            id_: str,
            aggregates: [str],
            validity: str,
            categories: [str],
    ):
        self.id_: str = id_
        self.aggregates: [str] = aggregates
        self.validity: str = validity
        self.categories: [str] = categories


class DFD(Cloneable):
    def __init__(
            self,
            data_types: [DataType],
            external_entities: [ExternalEntity],
            processes: [Process],
            data_stores: [DataStore],
            data_flows: [DataFlow],
    ):
        self.data_types: [DataType] = data_types
        self.external_entities: [ExternalEntity] = external_entities
        self.processes: [Process] = processes
        self.data_stores: [DataStore] = data_stores
        self.data_flows: [DataFlow] = data_flows
