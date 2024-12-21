# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# source: modelx.proto
# Protobuf Python Version: 5.26.1
"""Generated protocol buffer code."""
from google.protobuf import descriptor as _descriptor
from google.protobuf import descriptor_pool as _descriptor_pool
from google.protobuf import symbol_database as _symbol_database
from google.protobuf.internal import builder as _builder
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()




DESCRIPTOR = _descriptor_pool.Default().AddSerializedFile(b'\n\x0cmodelx.proto\x12\x05model\"\r\n\x0bPingRequest\"\x1f\n\tPongReply\x12\x12\n\x04\x63ode\x18\x01 \x01(\x05R\x04\x63ode\"F\n\x10\x45mbeddingRequest\x12\x12\n\x04text\x18\x01 \x01(\tR\x04text\x12\x1e\n\nmodel_name\x18\x02 \x01(\tR\nmodel_name\".\n\x0e\x45mbeddingReply\x12\x1c\n\tembedding\x18\x01 \x01(\tR\tembedding\"T\n\x14\x45mbeddingListRequest\x12\x1c\n\ttext_list\x18\x01 \x03(\tR\ttext_list\x12\x1e\n\nmodel_name\x18\x02 \x01(\tR\nmodel_name\"<\n\x12\x45mbeddingListReply\x12&\n\x0e\x65mbedding_list\x18\x01 \x03(\tR\x0e\x65mbedding_list\"y\n\x11SimilarityRequest\x12 \n\x0bsource_text\x18\x01 \x01(\tR\x0bsource_text\x12\"\n\x0ctarget_texts\x18\x02 \x03(\tR\x0ctarget_texts\x12\x1e\n\nmodel_name\x18\x03 \x01(\tR\nmodel_name\")\n\x0fSimilarityReply\x12\x16\n\x06scores\x18\x01 \x03(\x01R\x06scores\"=\n\x0f\x43lassifyRequest\x12\x12\n\x04text\x18\x01 \x01(\tR\x04text\x12\x16\n\x06labels\x18\x02 \x03(\tR\x06labels\"?\n\rClassifyReply\x12\x16\n\x06labels\x18\x01 \x03(\tR\x06labels\x12\x16\n\x06scores\x18\x02 \x03(\x01R\x06scores\",\n\x16\x45xtractKeywordsRequest\x12\x12\n\x04text\x18\x01 \x01(\tR\x04text\"2\n\x14\x45xtractKeywordsReply\x12\x1a\n\x08keywords\x18\x01 \x03(\tR\x08keywords2\xe4\x02\n\x06Modelx\x12.\n\x04Ping\x12\x12.model.PingRequest\x1a\x10.model.PongReply\"\x00\x12@\n\x0cGenEmbedding\x12\x17.model.EmbeddingRequest\x1a\x15.model.EmbeddingReply\"\x00\x12L\n\x10GenEmbeddingList\x12\x1b.model.EmbeddingListRequest\x1a\x19.model.EmbeddingListReply\"\x00\x12I\n\x13\x43\x61lcSimilarityScore\x12\x18.model.SimilarityRequest\x1a\x16.model.SimilarityReply\"\x00\x12O\n\x0f\x45xtractKeywords\x12\x1d.model.ExtractKeywordsRequest\x1a\x1b.model.ExtractKeywordsReply\"\x00\x42\x0bZ\t./;modelxb\x06proto3')

_globals = globals()
_builder.BuildMessageAndEnumDescriptors(DESCRIPTOR, _globals)
_builder.BuildTopDescriptorsAndMessages(DESCRIPTOR, 'modelx_pb2', _globals)
if not _descriptor._USE_C_DESCRIPTORS:
  _globals['DESCRIPTOR']._loaded_options = None
  _globals['DESCRIPTOR']._serialized_options = b'Z\t./;modelx'
  _globals['_PINGREQUEST']._serialized_start=23
  _globals['_PINGREQUEST']._serialized_end=36
  _globals['_PONGREPLY']._serialized_start=38
  _globals['_PONGREPLY']._serialized_end=69
  _globals['_EMBEDDINGREQUEST']._serialized_start=71
  _globals['_EMBEDDINGREQUEST']._serialized_end=141
  _globals['_EMBEDDINGREPLY']._serialized_start=143
  _globals['_EMBEDDINGREPLY']._serialized_end=189
  _globals['_EMBEDDINGLISTREQUEST']._serialized_start=191
  _globals['_EMBEDDINGLISTREQUEST']._serialized_end=275
  _globals['_EMBEDDINGLISTREPLY']._serialized_start=277
  _globals['_EMBEDDINGLISTREPLY']._serialized_end=337
  _globals['_SIMILARITYREQUEST']._serialized_start=339
  _globals['_SIMILARITYREQUEST']._serialized_end=460
  _globals['_SIMILARITYREPLY']._serialized_start=462
  _globals['_SIMILARITYREPLY']._serialized_end=503
  _globals['_CLASSIFYREQUEST']._serialized_start=505
  _globals['_CLASSIFYREQUEST']._serialized_end=566
  _globals['_CLASSIFYREPLY']._serialized_start=568
  _globals['_CLASSIFYREPLY']._serialized_end=631
  _globals['_EXTRACTKEYWORDSREQUEST']._serialized_start=633
  _globals['_EXTRACTKEYWORDSREQUEST']._serialized_end=677
  _globals['_EXTRACTKEYWORDSREPLY']._serialized_start=679
  _globals['_EXTRACTKEYWORDSREPLY']._serialized_end=729
  _globals['_MODELX']._serialized_start=732
  _globals['_MODELX']._serialized_end=1088
# @@protoc_insertion_point(module_scope)