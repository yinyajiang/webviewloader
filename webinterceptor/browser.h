#pragma once

#include <QMainWindow>
#include <QWebEngineView>
#include <QWebEnginePage>
#include <QLineEdit>
#include <QLabel>

class WebInterceptor;

class Browser : public QMainWindow {
    Q_OBJECT
public:
    Browser(const QString& url, const QString& title, const QString& ua,
            int width, int height, const QString& banner,
            const QString& bannerColor, bool showAddress, const QString& winColor, const QString& bannerFontColor, bool isforever);
    ~Browser();

protected:
    void closeEvent(QCloseEvent* event) override;

private slots:
    void loadUrl();
    void urlChanged(const QUrl& url);

private:
    QWebEngineView* m_webView = nullptr;
    QWebEnginePage* m_page = nullptr;
    WebInterceptor* m_profile = nullptr;
    QLineEdit* m_urlEdit = nullptr;
};
