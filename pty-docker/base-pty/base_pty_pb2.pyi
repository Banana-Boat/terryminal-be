from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Optional as _Optional

DESCRIPTOR: _descriptor.FileDescriptor

class RunCmdRequest(_message.Message):
    __slots__ = ["cmd"]
    CMD_FIELD_NUMBER: _ClassVar[int]
    cmd: str
    def __init__(self, cmd: _Optional[str] = ...) -> None: ...

class RunCmdResponse(_message.Message):
    __slots__ = ["result"]
    RESULT_FIELD_NUMBER: _ClassVar[int]
    result: str
    def __init__(self, result: _Optional[str] = ...) -> None: ...
