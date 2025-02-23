from threading import RLock

from src.models.session import Session


class SessionRepository:
    def __init__(self) -> None:
        self._sessions: dict[str, Session] = {}
        self._lock = RLock()

    def get_all_session(self) -> list[Session]:
        return list(self._sessions.values())

    def get_session(self, session_id: str) -> Session | None:
        return self._sessions.get(session_id)

    def create_session(self, session: Session) -> Session:
        with self._lock:
            self._sessions[session.id] = session
            return session

    def update_session(self, session: Session) -> Session:
        with self._lock:
            self._sessions[session.id] = session
            return session

    def delete_session(self, session_id: str) -> None:
        with self._lock:
            del self._sessions[session_id]
