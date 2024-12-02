REM Check if the virtual environment exists
IF NOT EXIST "pyinstaller_venv" (
    REM Create virtual environment
    python -m venv pyinstaller_venv
)

REM Activate virtual environment
CALL pyinstaller_venv\Scripts\activate

pip install -r requirements.txt
pip install pyinstaller

REM Build
python build_pyinstaller.py %*

REM 查找dist目录下的.exe文件
set exe_file=%dist_dir%\*.exe

REM 计算exe的md5
echo %exe_file%: > %exe_file%.md5
md5sum %exe_file% >> %exe_file%.md5
