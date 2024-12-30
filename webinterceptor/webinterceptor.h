#pragma once

#include <QWebEngineProfile>
#include <QWebEngineUrlRequestInterceptor>
#include <QNetworkCookie>
#include <QWebEngineCookieStore>
#include <QWebEngineView>
#include <QSet>
#include "opt.h"

class UrlRequestInterceptor : public QWebEngineUrlRequestInterceptor {
    Q_OBJECT
public:
    UrlRequestInterceptor(QWebEngineProfile* profile, QWebEngineView* webView, Options opt);
    void interceptRequest(QWebEngineUrlRequestInfo& info) override;

private:
    bool isPlayable(const QString& url) const;
    QWebEngineProfile* m_profile;
    QList<QNetworkCookie> m_cookies;
    QString m_htmlTitle;
    QWebEngineView* m_webView;
    QSet<QString> m_allUrls;
    Options m_opt;
};

class WebInterceptor : public QWebEngineProfile {
    Q_OBJECT
public:
    WebInterceptor(QWebEngineView* webView, Options opt);

private:
    UrlRequestInterceptor* m_interceptor;
}; 
