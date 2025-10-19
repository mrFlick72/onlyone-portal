import threading
from app.revenue.domain.repository import RevenueRepository
from app.revenue.domain.service import SaveRevenue


class RevenueConfigurationProvider:

    _instance = None
    _lock = threading.Lock()

    @staticmethod
    def __get_instance():
        if RevenueConfigurationProvider._instance is None:
            with RevenueConfigurationProvider._lock:
                if RevenueConfigurationProvider._instance is None:
                    RevenueConfigurationProvider._instance = RevenueConfigurationProvider()
        return RevenueConfigurationProvider._instance
    
    @staticmethod
    def get_revenue_repository() -> RevenueRepository:
        pass
    
    @staticmethod
    def get_save_revenue_service() -> SaveRevenue:
        return SaveRevenue()
