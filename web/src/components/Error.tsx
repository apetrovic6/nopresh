export default function Error({ error }: { error: string }) {
  return (
    <div>
      <p className="text-red-400">{error}</p>
    </div>
  )
}
