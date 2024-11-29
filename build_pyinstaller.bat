REM Check if the virtual environment exists
IF NOT EXIST "venv" (
    REM Create virtual environment
    python -m venv pyinstaller_venv
)

REM Activate virtual environment
CALL pyinstaller_venv\Scripts\activate

pip install -r requirements.txt
pip install pyinstaller

REM Build
python build_pyinstaller.py %*
