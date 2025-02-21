import threading

from src.app.grpc import start_grpc_service
from src.config import load_yaml_config
from src.interrupt import setup_interrupt_handler
from src.logging import setup_logging


def main() -> None:
    """Main entry point for the application"""

    config = load_yaml_config("config.yml")
    logger = setup_logging(config.log)

    logger.info("Starting application")

    interrupt_event = setup_interrupt_handler()

    grpc_service_thread = threading.Thread(
        target=start_grpc_service, args=(config.grpc, interrupt_event)
    )
    grpc_service_thread.start()

    # Wait for the interrupt event to be set
    interrupt_event.wait()

    grpc_service_thread.join()

    logger.info("Application shutdown complete")


if __name__ == "__main__":
    main()
