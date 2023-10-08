FROM python:3.9-slim
WORKDIR /app
RUN apt install --no-install-recommends bash
RUN pip config set global.index-url https://mirrors.aliyun.com/pypi/simple/ && \
    pip install --no-cache-dir python-dotenv grpcio protobuf
RUN rm -rf /var/lib/apt/lists/*
COPY . .

EXPOSE 3220
CMD [ "python", "run.py"]