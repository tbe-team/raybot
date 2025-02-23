# type: ignore
import logging
import threading

from apscheduler.executors.pool import ThreadPoolExecutor
from apscheduler.jobstores.memory import MemoryJobStore
from apscheduler.schedulers.background import BackgroundScheduler

from src.app.application import Application
from src.controllers.job.clean_up_session_job import CleanUpSessionJobHandler

logger = logging.getLogger(__name__)


def start_background_job_service(
    app: Application, interupt_event: threading.Event
) -> None:
    scheduler = BackgroundScheduler(
        jobstores={
            "default": MemoryJobStore(),
        },
        executors={
            "default": ThreadPoolExecutor(10),
        },
        job_default={},
    )

    register_job_handlers(app, scheduler)

    logger.info("Starting background job service")
    scheduler.start()

    # Block until the interrupt event is set
    interupt_event.wait()

    logger.info("Shutting down background job service...")
    scheduler.shutdown(wait=True)
    logger.info("Background job service shutdown complete")


def register_job_handlers(app: Application, scheduler: BackgroundScheduler) -> None:
    clean_up_session_handler = CleanUpSessionJobHandler(
        app.services.background_job_service
    )
    scheduler.add_job(
        clean_up_session_handler.handle,
        "interval",
        minutes=1,
        id="clean_up_session",
    )
