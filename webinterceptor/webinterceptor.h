#pragma once

#include <QWebEngineProfile>
#include <QWebEngineUrlRequestInterceptor>
#include <QNetworkCookie>
#include <QWebEngineCookieStore>

class UrlRequestInterceptor : public QWebEngineUrlRequestInterceptor {
    Q_OBJECT
public:
    UrlRequestInterceptor(QWebEngineProfile* profile, QObject* parent = nullptr);
    void interceptRequest(QWebEngineUrlRequestInfo& info) override;

private:
    bool isPlayable(const QString& url) const;
    QWebEngineProfile* m_profile;
    QList<QNetworkCookie> m_cookies;
};

class WebInterceptor : public QWebEngineProfile {
    Q_OBJECT
public:
    WebInterceptor(const QString& ua, QObject* parent = nullptr);

private:
    UrlRequestInterceptor* m_interceptor;
}; 