from src.config import AppConfig, load_yaml_config
from src.log import setup_log
from src.repositories.session_repository import SessionRepository
from src.services.auth_service import AuthService
from src.services.background_job_service import BackgroundJobService


class Application:
    """Application class for the Raybot project

    This class is responsible for initializing the application and providing access to the
    config, repositories, and services.

    """

    def __init__(self) -> None:
        config = load_yaml_config("config.yml")
        logger = setup_log(config.log)

        logger.info("Starting application")

        self._config = config
        self._repositories = Repositories()
        self._services = Services(config, self._repositories)

    @property
    def config(self) -> AppConfig:
        return self._config

    @property
    def services(self) -> "Services":
        return self._services


class Repositories:
    """Repositories class for the Raybot project

    This class is responsible for initializing the repositories and providing access to the
    repositories.

    """

    def __init__(self) -> None:
        self._session_repository = SessionRepository()

    @property
    def session_repository(self) -> SessionRepository:
        return self._session_repository


class Services:
    """Services class for the Raybot project

    This class is responsible for initializing the services and providing access to the
    services.

    """

    def __init__(self, config: AppConfig, repositories: Repositories) -> None:
        self._auth_service = AuthService(config, repositories.session_repository)
        self._background_job_service = BackgroundJobService(
            config, repositories.session_repository
        )

    @property
    def auth_service(self) -> AuthService:
        return self._auth_service

    @property
    def background_job_service(self) -> BackgroundJobService:
        return self._background_job_service
