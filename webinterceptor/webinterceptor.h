#pragma once

#include <QWebEngineProfile>
#include <QWebEngineUrlRequestInterceptor>
#include <QNetworkCookie>
#include <QWebEngineCookieStore>
#include <QWebEngineView>

class UrlRequestInterceptor : public QWebEngineUrlRequestInterceptor {
    Q_OBJECT
public:
    UrlRequestInterceptor(QWebEngineProfile* profile, QWebEngineView* webView, QObject* parent = nullptr, bool isforever = false);
    void interceptRequest(QWebEngineUrlRequestInfo& info) override;

private:
    bool isPlayable(const QString& url) const;
    QWebEngineProfile* m_profile;
    QList<QNetworkCookie> m_cookies;
    QString m_htmlTitle;
    QWebEngineView* m_webView;
    bool m_forever;
};

class WebInterceptor : public QWebEngineProfile {
    Q_OBJECT
public:
    WebInterceptor(const QString& ua, QWebEngineView* webView, QObject* parent = nullptr, bool isforever = false);

private:
    UrlRequestInterceptor* m_interceptor;
}; 
