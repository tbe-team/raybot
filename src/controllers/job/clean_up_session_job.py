import logging

from src.services.background_job_service import BackgroundJobService

logger = logging.getLogger(__name__)


class CleanUpSessionJobHandler:
    def __init__(self, background_job_service: BackgroundJobService) -> None:
        self._background_job_service = background_job_service

    def handle(self) -> None:
        logger.debug("Handle clean up expire session")
        try:
            self._background_job_service.clean_up_expire_session()
        except Exception as e:
            logger.error("Error clean up expire session: %s", e)
