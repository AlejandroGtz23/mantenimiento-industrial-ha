FROM nginx:1.27-alpine
COPY index.html.template /opt/balance/index.html.template
COPY entrypoint.sh /opt/balance/entrypoint.sh
RUN chmod +x /opt/balance/entrypoint.sh
ENTRYPOINT ["/opt/balance/entrypoint.sh"]
