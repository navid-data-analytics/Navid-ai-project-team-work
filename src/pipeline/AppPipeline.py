import logging
from src.components import Broadcaster, Joiner


logger = logging.getLogger('root')


class AppPipeline:
    def __init__(self, config, appID, AidServiceConnection, metric=None):
        """
        Responsible for processing the data for one appID.

                    /- trend detection --- decision maker ----------i
                   /                                                I
        Broadcaster --- fluctuation detection -- decision maker?-- Joiner
                  I                                                /
                  I___ trend prediction_______ decision maker_____/

        Arguments:
        - config dict, containing appropriate config for the type of pipeline
          with appIDs as keys for app specific configs
        - appID: int, appID
        - AidServiceConnection: connection to AID-E
        """

        self._appID = appID
        self._metric = metric
        self._AidServiceConnection = AidServiceConnection
        # init components
        self._entrypoint = Broadcaster()
        self._exitpoint = Joiner(
            process=self._joiner_process)

        self._create_components(config)
        self._link_components()

        # define input and output
        self.input = self._entrypoint.input
        self.output = self._exitpoint.output

    def _create_components(self, config):
        raise 'Not Implemented!'

    def _joiner_process(self):
        raise 'Not Implemented!'

    def _link_components(self):
        raise 'Not Implemented!'

    @property
    def appID(self):
        return self._appID

    @property
    def metric(self):
        return self._metric
