import os
import sys
import json

# Flask
from flask import Flask, request, jsonify
from gevent.pywsgi import WSGIServer
import tensorflow_hub as hub
import tensorflow as tf

# TensorFlow and tf.keras

from tensorflow.keras.applications.imagenet_utils import preprocess_input, decode_predictions
from tensorflow.keras.preprocessing import image

# Some utilites
import numpy as np
import io
from util import *

# Declare a flask app
app = Flask(__name__)

hub_module = hub.load("magenta_arbitrary-image-stylization-v1-256_2")

print('Model loaded. Check http://127.0.0.1:3000/')

# Model saved with Keras model.save()
MODEL_PATH = 'models/your_model.h5'

# Load your own trained model
# model = load_model(MODEL_PATH)
# model._make_predict_function()          # Necessary
# print('Model loaded. Start serving...')


# 风格迁移函数
def style_transfer(content_image, style_image):
    # 将图像转换为模型所需的格式
    content = tf.image.convert_image_dtype(content_image, tf.float32)
    style = tf.image.convert_image_dtype(style_image, tf.float32)
    content = tf.expand_dims(content, axis=0)
    style = tf.expand_dims(style, axis=0)

    # 进行风格迁移
    stylized_image = hub_module(tf.constant(content), tf.constant(style))[0]
    stylized_image = tf.clip_by_value(stylized_image, 0, 1)

    # 将结果转换回PIL图像
    stylized_image = tf.image.convert_image_dtype(stylized_image, tf.uint8)
    # stylized_image = Image.fromarray(stylized_image.numpy())

    return stylized_image.numpy()[0]


@app.route('/style', methods=['POST'])
def transfer():
    # Get the image from post request
    data = json.loads(request.data)
    content = data['content']
    style = data['style']

    content_img = base64_to_pil(content)
    style_img = base64_to_pil(style)

    stylized_img = style_transfer(content_img, style_img)
    image = Image.fromarray(stylized_img)

    # 创建一个字节流对象
    buffer = io.BytesIO()
    image.save(buffer, format='JPEG')

    # 将图像字节流编码为Base64字符串
    base64_image = base64.b64encode(buffer.getvalue()).decode('utf-8')
    image.save("saved_image.jpg")

    # 返回Base64字符串作为响应
    return base64_image


if __name__ == '__main__':
    # app.run(port=5002, threaded=False)

    # Serve the app with gevent
    http_server = WSGIServer(('0.0.0.0', 3000), app)
    http_server.serve_forever()