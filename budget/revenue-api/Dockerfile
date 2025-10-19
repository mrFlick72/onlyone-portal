FROM amazoncorretto:21

COPY lib /usr/local/onlyone-portal/

VOLUME /var/log/onlyone-portal

WORKDIR /usr/local/onlyone-portal/

CMD ["java", "-jar", "application.jar"]