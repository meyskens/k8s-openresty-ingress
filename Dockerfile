ARG ARCH
# Build go binary
FROM golang AS gobuild

COPY ./ /go/src/github.com/meyskens/k8s-openresty-ingress
WORKDIR /go/src/github.com/meyskens/k8s-openresty-ingress/controller

ARG GOARCH
ARG GOARM

RUN GOARCH=${GOARCH} GOARM=${GOARM} go build ./

# Set up deinitive image
ARG ARCH
FROM multiarch/debian-debootstrap:${ARCH}-jessie

RUN echo "deb http://httpredir.debian.org/debian jessie-backports main" >>/etc/apt/sources.list

RUN apt-get update && apt-get -y upgrade
RUN apt-get install -y -t jessie-backports \
    curl perl make build-essential procps \
    libreadline-dev libncurses5-dev \
    libpcre3-dev libssl-dev openssl \
    luarocks unzip zlib1g-dev zlibc git

ENV OPENRESTY_VERSION 1.13.6.1
ENV OPENRESTY_PREFIX /opt/openresty
ENV NGINX_PREFIX /opt/openresty/nginx
ENV VAR_PREFIX /var/nginx

RUN cd /tmp \
 && curl -sSL http://openresty.org/download/openresty-${OPENRESTY_VERSION}.tar.gz | tar -xvz \
 && cd openresty-* \
 && readonly NPROC=$(grep -c ^processor /proc/cpuinfo 2>/dev/null || 1) \
 && ./configure \
    --prefix=$OPENRESTY_PREFIX \
    --http-client-body-temp-path=$VAR_PREFIX/client_body_temp \
    --http-proxy-temp-path=$VAR_PREFIX/proxy_temp \
    --http-log-path=$VAR_PREFIX/access.log \
    --conf-path=/etc/nginx/nginx.conf \
    --error-log-path=$VAR_PREFIX/error.log \
    --pid-path=$VAR_PREFIX/nginx.pid \
    --lock-path=$VAR_PREFIX/nginx.lock \
    --with-luajit \
    --with-pcre-jit \
    --with-ipv6 \
    --with-threads \
    --with-http_v2_module \
    --with-http_ssl_module \
    --without-http_ssi_module \
    --without-http_userid_module \
    --without-http_uwsgi_module \
    --without-http_scgi_module \
    -j${NPROC} \
 && make -j${NPROC} \
 && make install \
 && ln -sf $NGINX_PREFIX/sbin/nginx /usr/local/bin/nginx \
 && ln -sf $NGINX_PREFIX/sbin/nginx /usr/local/bin/openresty \
 && ln -sf $OPENRESTY_PREFIX/bin/resty /usr/local/bin/resty \
 && ln -sf $OPENRESTY_PREFIX/luajit/bin/luajit-* $OPENRESTY_PREFIX/luajit/bin/lua \
 && ln -sf $OPENRESTY_PREFIX/luajit/bin/luajit-* /usr/local/bin/lua \
 && rm -rf /tmp/ngx_openresty*

COPY --from=gobuild /go/src/github.com/meyskens/k8s-openresty-ingress/controller/controller /usr/local/bin/controller

COPY ./config/default/ /etc/nginx/

COPY ./template/ingress.tpl /etc/nginx/ingress.tpl
ENV OPENRESTY_TEMPLATEPATH=/etc/nginx/ingress.tpl
ENV OPENRESTY_INGRESSATH=/etc/nginx/sites/

EXPOSE 80
EXPOSE 443
CMD controller