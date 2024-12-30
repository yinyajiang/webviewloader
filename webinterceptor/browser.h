#pragma once

#include <QMainWindow>
#include <QWebEngineView>
#include <QWebEnginePage>
#include <QLineEdit>
#include <QLabel>
#include "opt.h"

class WebInterceptor;

class Browser : public QMainWindow {
    Q_OBJECT
public:
    Browser(Options opt);
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
