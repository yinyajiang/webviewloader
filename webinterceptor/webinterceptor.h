#pragma once

#include <QWebEngineProfile>
#include <QWebEngineUrlRequestInterceptor>

class UrlRequestInterceptor : public QWebEngineUrlRequestInterceptor {
    Q_OBJECT
public:
    UrlRequestInterceptor(QObject* parent = nullptr);
    void interceptRequest(QWebEngineUrlRequestInfo& info) override;

private:
    bool isPlayable(const QString& url) const;
};

class WebInterceptor : public QWebEngineProfile {
    Q_OBJECT
public:
    WebInterceptor(const QString& ua, QObject* parent = nullptr);

private:
    UrlRequestInterceptor* interceptor;
}; 