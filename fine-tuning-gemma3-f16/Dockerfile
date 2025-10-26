FROM unsloth/unsloth:stable

# Install libcurl4-openssl-dev which is required by llama.cpp to save GGUF files and
# update the unsloth libraries to bring in some fixes related to saving to GGUF files
USER root
RUN apt update && \
    apt install libcurl4-openssl-dev -y && \
    rm -rf /var/lib/apt/lists/*
RUN pip install --upgrade unsloth-zoo && \
    pip install --upgrade unsloth

WORKDIR /usr/local/app

COPY --link finetune.py ./
ENTRYPOINT [ ]
CMD ["python", "finetune.py"]