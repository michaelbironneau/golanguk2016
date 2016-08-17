from werkzeug.wrappers import Request, Response
from werkzeug.serving import run_simple
import sys
from jsonrpc import JSONRPCResponseManager, dispatcher


@dispatcher.add_method
def fit(**kwargs):
    return {"TrainedModel": {"Label": kwargs["Labels"][0]}}

@dispatcher.add_method
def predict(**kwargs):
    return {"PredictedLabel": kwargs["TrainedModel"]["Label"]}


@Request.application
def application(request):
    response = JSONRPCResponseManager.handle(
        request.data, dispatcher)
    return Response(response.json, mimetype='application/json')


if __name__ == '__main__':
    run_simple('localhost', int(sys.argv[1]), application)