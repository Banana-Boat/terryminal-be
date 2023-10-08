from concurrent import futures
import os
import grpc
from dotenv import load_dotenv

from server import BasePtyServicer
from base_pty_pb2_grpc import add_BasePtyServicer_to_server


def serve(host, port):
    server = grpc.server(futures.ThreadPoolExecutor())
    add_BasePtyServicer_to_server(
        BasePtyServicer(), server
    )
    server.add_insecure_port(f"{host}:{port}")
    server.start()
    print(f"Server running at: {host}:{port}...")
    server.wait_for_termination()


if __name__ == "__main__":
    load_dotenv()
    host = os.getenv('SERVICE_HOST')
    port = os.getenv('SERVICE_PORT')

    serve(host, port)
