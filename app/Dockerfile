FROM ubuntu
RUN apt update -y
RUN apt install net-tools dnsutils nmap -y
RUN apt install curl git wget -y
RUN wget https://go.dev/dl/go1.24.0.linux-amd64.tar.gz && rm -rf /usr/local/go && tar -C /usr/local -xzf go1.24.0.linux-amd64.tar.gz
RUN export PATH=$PATH:/usr/local/go/bin && git clone -b v3.4.35 https://github.com/etcd-io/etcd.git && cd etcd && ./build && export PATH="$PATH:`pwd`/bin"
RUN curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl" && chmod +x kubectl && mv kubectl /usr/local/bin/
RUN apt install telnet vim -y
RUN tar -C /usr/local -xzf go1.24.0.linux-amd64.tar.gz && export PATH=$PATH:/usr/local/go/bin

WORKDIR /opt

COPY app-linux /usr/local/bin/

RUN chmod ugo+x /usr/local/bin/app-linux

# command: ["sleep", "infinity"]
CMD [app-linux]
