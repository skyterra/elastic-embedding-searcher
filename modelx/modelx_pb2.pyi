from google.protobuf.internal import containers as _containers
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Iterable as _Iterable, Optional as _Optional

DESCRIPTOR: _descriptor.FileDescriptor

class PingRequest(_message.Message):
    __slots__ = ()
    def __init__(self) -> None: ...

class PongReply(_message.Message):
    __slots__ = ("code",)
    CODE_FIELD_NUMBER: _ClassVar[int]
    code: int
    def __init__(self, code: _Optional[int] = ...) -> None: ...

class EmbeddingRequest(_message.Message):
    __slots__ = ("text", "model_name")
    TEXT_FIELD_NUMBER: _ClassVar[int]
    MODEL_NAME_FIELD_NUMBER: _ClassVar[int]
    text: str
    model_name: str
    def __init__(self, text: _Optional[str] = ..., model_name: _Optional[str] = ...) -> None: ...

class EmbeddingReply(_message.Message):
    __slots__ = ("embedding",)
    EMBEDDING_FIELD_NUMBER: _ClassVar[int]
    embedding: str
    def __init__(self, embedding: _Optional[str] = ...) -> None: ...

class EmbeddingListRequest(_message.Message):
    __slots__ = ("text_list", "model_name")
    TEXT_LIST_FIELD_NUMBER: _ClassVar[int]
    MODEL_NAME_FIELD_NUMBER: _ClassVar[int]
    text_list: _containers.RepeatedScalarFieldContainer[str]
    model_name: str
    def __init__(self, text_list: _Optional[_Iterable[str]] = ..., model_name: _Optional[str] = ...) -> None: ...

class EmbeddingListReply(_message.Message):
    __slots__ = ("embedding_list",)
    EMBEDDING_LIST_FIELD_NUMBER: _ClassVar[int]
    embedding_list: _containers.RepeatedScalarFieldContainer[str]
    def __init__(self, embedding_list: _Optional[_Iterable[str]] = ...) -> None: ...

class SimilarityRequest(_message.Message):
    __slots__ = ("source_text", "target_texts", "model_name")
    SOURCE_TEXT_FIELD_NUMBER: _ClassVar[int]
    TARGET_TEXTS_FIELD_NUMBER: _ClassVar[int]
    MODEL_NAME_FIELD_NUMBER: _ClassVar[int]
    source_text: str
    target_texts: _containers.RepeatedScalarFieldContainer[str]
    model_name: str
    def __init__(self, source_text: _Optional[str] = ..., target_texts: _Optional[_Iterable[str]] = ..., model_name: _Optional[str] = ...) -> None: ...

class SimilarityReply(_message.Message):
    __slots__ = ("scores",)
    SCORES_FIELD_NUMBER: _ClassVar[int]
    scores: _containers.RepeatedScalarFieldContainer[float]
    def __init__(self, scores: _Optional[_Iterable[float]] = ...) -> None: ...

class ClassifyRequest(_message.Message):
    __slots__ = ("text", "labels")
    TEXT_FIELD_NUMBER: _ClassVar[int]
    LABELS_FIELD_NUMBER: _ClassVar[int]
    text: str
    labels: _containers.RepeatedScalarFieldContainer[str]
    def __init__(self, text: _Optional[str] = ..., labels: _Optional[_Iterable[str]] = ...) -> None: ...

class ClassifyReply(_message.Message):
    __slots__ = ("labels", "scores")
    LABELS_FIELD_NUMBER: _ClassVar[int]
    SCORES_FIELD_NUMBER: _ClassVar[int]
    labels: _containers.RepeatedScalarFieldContainer[str]
    scores: _containers.RepeatedScalarFieldContainer[float]
    def __init__(self, labels: _Optional[_Iterable[str]] = ..., scores: _Optional[_Iterable[float]] = ...) -> None: ...

class ExtractKeywordsRequest(_message.Message):
    __slots__ = ("text",)
    TEXT_FIELD_NUMBER: _ClassVar[int]
    text: str
    def __init__(self, text: _Optional[str] = ...) -> None: ...

class ExtractKeywordsReply(_message.Message):
    __slots__ = ("keywords",)
    KEYWORDS_FIELD_NUMBER: _ClassVar[int]
    keywords: _containers.RepeatedScalarFieldContainer[str]
    def __init__(self, keywords: _Optional[_Iterable[str]] = ...) -> None: ...
