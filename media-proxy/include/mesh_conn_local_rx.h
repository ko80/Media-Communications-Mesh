#ifndef MESH_CONN_LOCAL_RX_H
#define MESH_CONN_LOCAL_RX_H

#include "mesh_conn_local.h"
#include "shm_memif.h"

namespace mesh {

namespace connection {

class LocalRx : public Local {
public:
    LocalRx(memif_ops_t *ops);
    ~LocalRx() override;

    memif_ops_t ops;

private:
    Result on_establish(context::Context& ctx) override;
    Result on_shutdown(context::Context& ctx) override;
    void on_delete(context::Context& ctx) override;
};

} // namespace connection

} // namespace mesh

#endif // MESH_CONN_LOCAL_RX_H
