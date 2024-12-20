#include "webinterceptor.h"
#include <QWebEngineSettings>
#include <QJsonObject>
#include <QJsonDocument>
#include <QCoreApplication>
#include <iostream>
#include <QNetworkCookie>

UrlRequestInterceptor::UrlRequestInterceptor(QWebEngineProfile* profile, QObject* parent)
    : QWebEngineUrlRequestInterceptor(parent)
    , m_profile(profile) {

    connect(m_profile->cookieStore(), &QWebEngineCookieStore::cookieAdded,
            this, [this](const QNetworkCookie& cookie) {
                m_cookies.append(cookie);
            });
    connect(m_profile->cookieStore(), &QWebEngineCookieStore::cookieRemoved,
            this, [this](const QNetworkCookie& cookie) {
                m_cookies.removeAll(cookie);
            });

}

void UrlRequestInterceptor::interceptRequest(QWebEngineUrlRequestInfo& info) {
    QString url = info.requestUrl().toString();
    if (isPlayable(url)) {
        QJsonObject obj;
        obj["url"] = url;
        obj["method"] = QString::fromUtf8(info.requestMethod());
        
        // 添加 headers
        QJsonObject headers;
        QHash<QByteArray, QByteArray> httpHeaders = info.httpHeaders();
        for (auto it = httpHeaders.constBegin(); it != httpHeaders.constEnd(); ++it) {
            headers[QString::fromUtf8(it.key())] = QString::fromUtf8(it.value());
        }
        obj["headers"] = headers;
        
        // 添加收集到的 cookies
        QJsonObject cookies;
        for (const QNetworkCookie& cookie : m_cookies) {
            QJsonObject cookieObj;
            cookieObj["name"] = QString::fromUtf8(cookie.name());
            cookieObj["value"] = QString::fromUtf8(cookie.value());
            cookieObj["domain"] = cookie.domain();
            cookieObj["path"] = cookie.path();
            cookieObj["expires"] = cookie.expirationDate().toString();
            cookies[cookie.domain()] = cookieObj;
        }
        obj["cookies"] = cookies;

        QJsonDocument doc(obj);
        std::cout << doc.toJson(QJsonDocument::Compact).toStdString() << std::endl;
        std::cout.flush();
        QCoreApplication::exit(0);
    }
}

bool UrlRequestInterceptor::isPlayable(const QString& urlString) const {
    QUrl url(urlString);
    QString path = url.path().toLower();
    
    const QStringList playableExtensions = {
        ".m3u8", ".mp4",  ".m4v",  ".mov",  ".flv", ".webm", ".ts",      
        ".mkv", ".avi",  ".mpd", ".f4v", ".m4s", ".3gp"     
    };
    for (const QString& ext : playableExtensions) {
        if (path.endsWith(ext)) {
            return true;
        }
    }
    return false;
}

WebInterceptor::WebInterceptor(const QString& ua, QObject* parent)
    : QWebEngineProfile(parent) {
    m_interceptor = new UrlRequestInterceptor(this, this);
    setUrlRequestInterceptor(m_interceptor);
    
    if (!ua.isEmpty()) {
        setHttpUserAgent(ua);
    }
    
    settings()->setAttribute(QWebEngineSettings::PlaybackRequiresUserGesture, true);
} 
