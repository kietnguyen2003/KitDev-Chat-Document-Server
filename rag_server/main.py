try:
    from .app import app
except ImportError:  # pragma: no cover
    import sys
    from pathlib import Path

    current_dir = Path(__file__).resolve().parent
    project_root = current_dir.parent
    if str(project_root) not in sys.path:
        sys.path.insert(0, str(project_root))

    from rag_server.app import app


__all__ = ["app"]
