import threading

from src.app.application import Application
from src.app.grpc import start_grpc_service  # type: ignore[attr-defined]
from src.app.job import start_background_job_service  # type: ignore[attr-defined]
from src.interrupt import setup_interrupt_handler


def main() -> None:
    app = Application()

    interrupt_event = setup_interrupt_handler()

    job_service_thread = threading.Thread(
        target=start_background_job_service, args=(app, interrupt_event)
    )
    job_service_thread.start()

    grpc_service_thread = threading.Thread(
        target=start_grpc_service, args=(app, interrupt_event)
    )
    grpc_service_thread.start()

    # Wait for the interrupt event to be set
    interrupt_event.wait()

    grpc_service_thread.join()
    job_service_thread.join()


if __name__ == "__main__":
    main()
