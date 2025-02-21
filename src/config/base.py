from pydantic import BaseModel


class ConfigModel(BaseModel):
    """Base config model

    All config models should inherit from this class.

    Example:
    ```
    class MyConfig(ConfigModel):
        my_field: int
    ```
    """

    class Config:
        use_enum_values = True
        populate_by_name = True
