import os
import sys

# Flask
from flask import Flask, request, jsonify
from gevent.pywsgi import WSGIServer

# TensorFlow and tf.keras

from tensorflow.keras.applications.imagenet_utils import preprocess_input, decode_predictions
from tensorflow.keras.preprocessing import image
from tensorflow.keras.activations import softmax

# Some utilites
import numpy as np
from util import *

# Declare a flask app
app = Flask(__name__)

from tensorflow.keras.applications.mobilenet_v2 import MobileNetV2

model = MobileNetV2(weights='imagenet')

print('Model loaded. Check http://127.0.0.1:4000/')

# Model saved with Keras model.save()
MODEL_PATH = 'models/your_model.h5'

# Load your own trained model
# model = load_model(MODEL_PATH)
# model._make_predict_function()          # Necessary
# print('Model loaded. Start serving...')


@app.route('/predict', methods=['POST'])
def predict():
    # Get the image from post request

    # img = base64_to_pil(request.json)
    img_str = request.get_data().decode('utf-8')[1:-1]
    img = base64_to_pil(img_str)

    # Make prediction
    preds = model_predict(img, model)

    # Process your result for human
    pred_proba = "{:.3f}".format(np.amax(preds))  # Max probability
    pred_class = decode_predictions(preds, top=1)  # ImageNet Decode

    result = str(pred_class[0][0][1])  # Convert to string
    result = result.replace('_', ' ').capitalize()

    # 生成热力图
    img_tmp = img.resize((224, 224))

    # Preprocessing the image
    img_arr = preprocess_input(image.img_to_array(img_tmp), mode='tf')
    model.layers[-1].activation = None
    target_layer_name = model.layers[-5].name
    heatmap = make_gradcam_heatmap(np.array([img_arr]), model,
                                   target_layer_name)
    save_and_display_gradcam(img, heatmap, alpha=1)
    model.layers[-1].activation = softmax

    # Serialize the result, you can add additional fields
    return jsonify(result=result, probability=pred_proba)


if __name__ == '__main__':
    # app.run(port=5002, threaded=False)

    # Serve the app with gevent
    http_server = WSGIServer(('0.0.0.0', 4000), app)
    http_server.serve_forever()