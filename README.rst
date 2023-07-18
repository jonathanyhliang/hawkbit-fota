Hawkbit FOTA Microservice
=========================

Overview
--------
``hawkbit-fota`` service consists of a frontend server which provides RESTful APIs to manage images,
distributions, and deployments to be deployed to ``hawkbit-fota`` clients. And a backend server which
implements `Hawkbit DDI <https://www.eclipse.org/hawkbit/apis/ddi_api/>`_ compliant RESTful APIs
so that ``hawkbit-fota`` clients could poll and launch FOTA processes. However, the backend server
implements part of Hawkbit DDI specified RESTful APIs that is just sufficient to host ``hawkbit-fota``
client devices with `Zephyr Hawkbit FOTA <https://github.com/jonathanyhliang/zephyr/tree/cc32xx-hawkbit-bringup/samples/subsys/mgmt/hawkbit>`_
sample application.

Building and Running
####################

The primary use of ``hawkbit-fota`` is run with ``demo-svc``. Refer to
`demo-svc <https://github.com/jonathanyhliang/demo-svc>`_ for the full picture of how things work.
