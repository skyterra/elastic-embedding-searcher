import yaml
import modelx_pb2
import modelx_pb2_grpc
from sklearn.metrics.pairwise import cosine_similarity
from sentence_transformers import SentenceTransformer


class Modelx(modelx_pb2_grpc.ModelxServicer):
    def __init__(self, model_path: str):
        self.model = SentenceTransformer(model_path)

    def Ping(self, request, context):
        return modelx_pb2.PongReply(code=0)

    def GenEmbedding(self, request, context):
        model = self.models.get(request.model_name)
        if model is None:
            return modelx_pb2.EmbeddingReply()

        embedding_vector = model.encode(request.text, show_progress_bar=False)
        embedding = ",".join(map(str, embedding_vector))
        return modelx_pb2.EmbeddingReply(embedding=embedding)

    def GenEmbeddingList(self, request, context):
        model = self.models.get(request.model_name)
        if model is None:
            return modelx_pb2.EmbeddingListReply()

        vector_list = model.encode(request.text_list, show_progress_bar=False)
        embedding_list = [",".join(map(str, vector)) for vector in vector_list]
        return modelx_pb2.EmbeddingListReply(embedding_list=embedding_list)

    def CalcSimilarityScore(self, request, context):
        model = self.models.get(request.model_name)
        if model is None:
            return modelx_pb2.SimilarityReply()

        source_text_embedding = model.encode([request.source_text], show_progress_bar=False)
        target_texts_embeddings = model.encode(request.target_texts, show_progress_bar=False)

        # mean pooling score.
        mean_pooling_scores = cosine_similarity(source_text_embedding, target_texts_embeddings)

        return modelx_pb2.SimilarityReply(scores=mean_pooling_scores[0])
