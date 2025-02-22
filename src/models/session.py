from dataclasses import dataclass
from datetime import datetime


@dataclass
class Session:
    id: str
    created_at: datetime
    last_seen_at: datetime
