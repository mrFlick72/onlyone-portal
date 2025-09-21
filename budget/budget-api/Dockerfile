FROM amazoncorretto:21

COPY target/application.jar /usr/local/onlyone-portal/

VOLUME /var/log/onlyone-portal

WORKDIR /usr/local/onlyone-portal/

CMD ["java", "-jar", "application.jar"]