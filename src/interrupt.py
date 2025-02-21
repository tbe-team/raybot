import signal
import threading
from typing import Any


def setup_interrupt_handler() -> threading.Event:
    """Setup an interrupt handler that will set an event when the application is interrupted.

    :return: An event that will be set when the application is interrupted.
    """
    interrupt_event = threading.Event()

    def handle_shutdown(_: int, __: Any) -> None:
        interrupt_event.set()

    signal.signal(signal.SIGINT, handle_shutdown)
    signal.signal(signal.SIGTERM, handle_shutdown)

    return interrupt_event
