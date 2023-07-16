FROM python:3.10

WORKDIR /app

# COPY requirements_no_version.txt requirements_no_version.txt

RUN pip config set global.index-url http://mirrors.aliyun.com/pypi/simple
RUN pip config set global.trusted-host mirrors.aliyun.com

RUN pip install flask
RUN pip install gevent
RUN pip install tensorflow
RUN pip install pillow
RUN pip install matplotlib
RUN pip install IPython


COPY util.py util.py
COPY classify_app.py classify_app.py


CMD ["python3","classify_app.py","--host=0.0.0.0"]