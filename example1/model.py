import sys 
import json

def fit(data, labels):
    """Fit the model given data and labels"""
    print('{"p": ' + labels[0] + '}')

def predict(data, model):
    """Return a prediction given input"""
    print('{"value": ' + model['p'] + '}')

if sys.argv[1] == 'fit':
    if len(sys.argv) != 4:
        print("{'error': 'incorrect number of arguments'}")
        sys.exit(1)
    else:
        d = json.loads(sys.argv[2])
        l = json.loads(sys.argv[3])
        fit(d,l)
elif sys.argv[1] == 'predict':
    if len(sys.argv) != 4:
        print("{'error': 'incorrect number of arguments'}")
        sys.exit(1)
    d = json.loads(sys.argv[2])
    m = json.loads(sys.argv[3])
    predict(data,  model)
sys.exit(0)