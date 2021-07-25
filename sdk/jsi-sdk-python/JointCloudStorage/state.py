from .bucket import _BASE


class State(_BASE):
    CATEGORY_STATE = "state"

    def get_storage_state(self):
        rs = self._do("GET", State.CATEGORY_STATE, "/storage")
        result = rs.json()
        return result

    def get_server_info(self):
        rs = self._do("GET", State.CATEGORY_STATE, "/server")
        result = rs.json()
        return result

    def get_task_info(self, task_id: str):
        task_id = task_id.strip()
        rs = self._do("GET", State.CATEGORY_STATE, "/server/" + task_id)
        result = rs.json()
        return result

    def get_storage_plan(self):
        rs = self._do("GET", State.CATEGORY_STATE, "/plan")
        result = rs.json()
        return result
