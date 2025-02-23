from typing import Protocol


class JobHandler(Protocol):
    def run(self) -> None: ...
