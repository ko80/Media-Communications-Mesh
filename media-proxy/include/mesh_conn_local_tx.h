#ifndef MESH_CONN_LOCAL_TX_H
#define MESH_CONN_LOCAL_TX_H

#include "mesh_conn_local.h"

namespace mesh {

namespace connection {

class LocalTx : public Local {

public:
    LocalTx();
    ~LocalTx() override;

private:
    Result on_establish(context::Context& ctx) override;
    Result on_receive(context::Context& ctx, void *ptr, uint32_t sz,
                      uint32_t& sent) override;
    Result on_shutdown(context::Context& ctx) override;
    void   on_delete(context::Context& ctx) override;

private:
};

} // namespace connection

} // namespace mesh

#endif // MESH_CONN_LOCAL_TX_H
