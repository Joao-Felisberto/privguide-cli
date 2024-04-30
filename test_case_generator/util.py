from abc import ABC
from typing import Self


class Cloneable(ABC):

    def clone(self, **kwargs) -> Self:
        params = vars(self)
        params = {key: kwargs.get(key, val) for key, val in params.items()}
        params = {
            k: params[k].clone() if isinstance(params[k], Cloneable) else
            params[k][:] if isinstance(params[k], list) else
            params[k]
            for k in params
        }

        return self.__class__(**params)
