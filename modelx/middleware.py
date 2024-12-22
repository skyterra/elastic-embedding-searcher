import grpc
import logging
import time
from grpc import ServerInterceptor


class LatencyInterceptor(ServerInterceptor):
    def intercept_service(self, continuation, handler_call_details):
        if handler_call_details.method.endswith('/Ping'):
            return continuation(handler_call_details)
        
        start_time = time.time_ns()
        request_id = dict(handler_call_details.invocation_metadata).get('request-id', 'unknown')
        
        response = continuation(handler_call_details)
        latency = (time.time_ns() - start_time)/1e6

        logging.info(f"{request_id}| {handler_call_details.method} | {latency:.3f} ms")
        return response
