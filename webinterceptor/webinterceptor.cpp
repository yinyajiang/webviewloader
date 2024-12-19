#include "webinterceptor.h"
#include <QWebEngineSettings>
#include <QJsonObject>
#include <QJsonDocument>
#include <QCoreApplication>
#include <iostream>

UrlRequestInterceptor::UrlRequestInterceptor(QObject* parent)
    : QWebEngineUrlRequestInterceptor(parent) {}

void UrlRequestInterceptor::interceptRequest(QWebEngineUrlRequestInfo& info) {
    QString url = info.requestUrl().toString();
    if (isPlayable(url)) {
        QJsonObject headers;
        try {
            QHash<QByteArray, QByteArray> httpHeaders = info.httpHeaders();
            for (auto it = httpHeaders.constBegin(); it != httpHeaders.constEnd(); ++it) {
                headers[QString::fromUtf8(it.key())] = QString::fromUtf8(it.value());
            }
        } catch (const std::exception& e) {
            qDebug() << "Error getting headers:" << e.what();
        }

        QJsonObject obj;
        obj["url"] = url;
        obj["method"] = QString::fromUtf8(info.requestMethod());
        obj["headers"] = headers;

        QJsonDocument doc(obj);
        std::cout << doc.toJson(QJsonDocument::Compact).toStdString() << std::endl;
        std::cout.flush();
        QCoreApplication::exit(0);
    }
}

bool UrlRequestInterceptor::isPlayable(const QString& url) const {
    return url.toLower().contains(".m3u8");
}

WebInterceptor::WebInterceptor(const QString& ua, QObject* parent)
    : QWebEngineProfile(parent) {
    interceptor = new UrlRequestInterceptor(this);
    setUrlRequestInterceptor(interceptor);
    
    if (!ua.isEmpty()) {
        setHttpUserAgent(ua);
    }
    
    settings()->setAttribute(QWebEngineSettings::PlaybackRequiresUserGesture, true);
} 
