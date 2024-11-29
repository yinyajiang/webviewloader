REM Check if the virtual environment exists
IF NOT EXIST "venv" (
    REM Create virtual environment
    python -m venv venv
)

REM Activate virtual environment
CALL venv\Scripts\activate

pip install -r requirements.txt
pip install pyinstaller

REM Build
python pyinstaller_build.py %*
