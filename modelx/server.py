import argparse
import signal
import sys
import logging
import os
import grpc
import modelx_pb2_grpc
from concurrent import futures
from modelx import Modelx
from middleware import LatencyInterceptor

DEFAULT_MAX_WORKERS = 5


def setup_logging():
    # config logging.
    logging_config = {
        "level": logging.INFO,
        "format": "%(asctime)s.%(msecs)03d|%(levelname)s|%(message)s",
        "datefmt": "%Y-%m-%d %H:%M:%S"
    }

    logging.basicConfig(**logging_config)


def serve(max_workers: int, is_uds: bool, model_path: str):
    logging.info("start modelx.")
    address = "unix:/tmp/grpc_unix_socket_modelx"
    if not is_uds:
        address = "[::]:50051"

    if is_uds and os.path.exists(address):
        os.remove(address)

    server = grpc.server(
        futures.ThreadPoolExecutor(max_workers=max_workers),
        interceptors=[LatencyInterceptor()],
        options=[
            ('grpc.max_receive_message_length', 10 << 20),  # 10 MB
            ('grpc.max_send_message_length', 10 << 20)      # 10 MB
        ]
    )

    modelx_pb2_grpc.add_ModelxServicer_to_server(Modelx(model_path), server)
    server.add_insecure_port(address)
    server.start()

    # notify parent process that I am ready.
    if is_uds:
        os.kill(os.getppid(), signal.SIGUSR1)
        logging.info("send SIGUSR1:" + os.getppid().__str__())

    logging.info("modelx start at:" + address)
    server.wait_for_termination()


def handle_exit(signum, frame):
    logging.info("receive sys signal" + signum + "frame" + frame)
    cleanup()


def cleanup():
    logging.info("cleanup.")
    logging.info("bye")
    sys.exit(0)


if __name__ == "__main__":
    # parse arguments
    parser = argparse.ArgumentParser(description="Modelx gGRPC Server")
    parser.add_argument("--uds", dest="is_uds", type=bool, default=False,
                        help=f"Enable unix domain socket")
    parser.add_argument("--max_workers", dest="max_workers", type=int, default=DEFAULT_MAX_WORKERS,
                        help=f"Set the maximum number of workers (default: {DEFAULT_MAX_WORKERS})")
    parser.add_argument("--model_path", dest="model_path", type=str, default="",
                        help=f"Set the model path")

    args = parser.parse_args()
    setup_logging()

    # register SIGINT and SIGTERM in order exit gracefully.
    signal.signal(signal.SIGINT | signal.SIGTERM, handle_exit)

    try:
        serve(args.max_workers, args.is_uds, args.model_path)
    except KeyboardInterrupt:
        logging.info("receive KeyboardInterrupt, exiting gracefully.")
        cleanup()
