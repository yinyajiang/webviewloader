#include "webinterceptor.h"
#include <QWebEngineSettings>
#include <QJsonObject>
#include <QJsonDocument>
#include <QCoreApplication>
#include <iostream>
#include <QNetworkCookie>
#include <QFile>

UrlRequestInterceptor::UrlRequestInterceptor(QWebEngineProfile* profile, QWebEngineView* webView, Options opt)
    : QWebEngineUrlRequestInterceptor(profile)
    , m_profile(profile)
    , m_webView(webView)
    , m_opt(opt){

    connect(m_profile->cookieStore(), &QWebEngineCookieStore::cookieAdded,
            this, [this](const QNetworkCookie& cookie) {
                m_cookies.append(cookie);
            });
    connect(m_profile->cookieStore(), &QWebEngineCookieStore::cookieRemoved,
            this, [this](const QNetworkCookie& cookie) {
                m_cookies.removeAll(cookie);
            });

    connect(m_webView, &QWebEngineView::titleChanged, this, [this](const QString& title) {
        this->m_htmlTitle = title;
    });

    if(!m_opt.dumpHtml.isEmpty()){
        connect(m_webView, &QWebEngineView::loadFinished, this, [this](bool ok) {
            if(!ok){
                return;
            }
            m_webView->page()->toHtml([this](const QString& html) {
                QFile file(m_opt.dumpHtml);
                if(file.open(QIODevice::WriteOnly | QIODevice::Text)){
                    file.write(html.toUtf8());
                    file.close();
                    QCoreApplication::exit(0);
                }
            });
        });
    }
}

void UrlRequestInterceptor::interceptRequest(QWebEngineUrlRequestInfo& info) {
    QString url = info.requestUrl().toString();
    if (isPlayable(url) && (m_allUrls.empty() || !m_allUrls.contains(url)) ) {
        QJsonObject obj;
        obj["url"] = url;
        obj["method"] = QString::fromUtf8(info.requestMethod());
        obj["title"] = m_htmlTitle;
        
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
        if(!m_opt.isforever){
            QCoreApplication::exit(0);
        }else{
            m_allUrls.insert(url);
        }
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


WebInterceptor::WebInterceptor(QWebEngineView* webView, Options opt)
    : QWebEngineProfile(webView) {
    m_interceptor = new UrlRequestInterceptor(this, webView, opt);
    setUrlRequestInterceptor(m_interceptor);
    setSpellCheckEnabled(false);
    
    if (!opt.ua.isEmpty()) {
        setHttpUserAgent(opt.ua);
    }

    settings()->setAttribute(QWebEngineSettings::PlaybackRequiresUserGesture, true);
} 
