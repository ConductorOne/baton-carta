FROM gcr.io/distroless/static-debian11:nonroot
ENTRYPOINT ["/baton-carta"]
COPY baton-carta /