# Generated by the gRPC Python protocol compiler plugin. DO NOT EDIT!
"""Client and server classes corresponding to protobuf-defined services."""
import grpc

import base_pty_pb2 as base__pty__pb2


class BasePtyStub(object):
    """
    Server: BasePty
    Client: Terminal Service
    """

    def __init__(self, channel):
        """Constructor.

        Args:
            channel: A grpc.Channel.
        """
        self.RunCmd = channel.stream_stream(
                '/BasePty/RunCmd',
                request_serializer=base__pty__pb2.RunCmdRequest.SerializeToString,
                response_deserializer=base__pty__pb2.RunCmdResponse.FromString,
                )


class BasePtyServicer(object):
    """
    Server: BasePty
    Client: Terminal Service
    """

    def RunCmd(self, request_iterator, context):
        """在BasePty中执行命令，返回结果
        """
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')


def add_BasePtyServicer_to_server(servicer, server):
    rpc_method_handlers = {
            'RunCmd': grpc.stream_stream_rpc_method_handler(
                    servicer.RunCmd,
                    request_deserializer=base__pty__pb2.RunCmdRequest.FromString,
                    response_serializer=base__pty__pb2.RunCmdResponse.SerializeToString,
            ),
    }
    generic_handler = grpc.method_handlers_generic_handler(
            'BasePty', rpc_method_handlers)
    server.add_generic_rpc_handlers((generic_handler,))


 # This class is part of an EXPERIMENTAL API.
class BasePty(object):
    """
    Server: BasePty
    Client: Terminal Service
    """

    @staticmethod
    def RunCmd(request_iterator,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.stream_stream(request_iterator, target, '/BasePty/RunCmd',
            base__pty__pb2.RunCmdRequest.SerializeToString,
            base__pty__pb2.RunCmdResponse.FromString,
            options, channel_credentials,
            insecure, call_credentials, compression, wait_for_ready, timeout, metadata)
